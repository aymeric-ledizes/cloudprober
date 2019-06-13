// Copyright 2017-2019 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Package prober provides a prober for running a set of probes.

Prober takes in a config proto which dictates what probes should be created
with what configuration, and manages the asynchronous fan-in/fan-out of the
metrics data from these probes.
*/
package prober

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/google/cloudprober/config"
	configpb "github.com/google/cloudprober/config/proto"
	"github.com/google/cloudprober/logger"
	"github.com/google/cloudprober/metrics"
	"github.com/google/cloudprober/probes"
	"github.com/google/cloudprober/servers"
	"github.com/google/cloudprober/surfacers"
	"github.com/google/cloudprober/sysvars"
	"github.com/google/cloudprober/targets/lameduck"
	rdsserver "github.com/google/cloudprober/targets/rds/server"
	"github.com/google/cloudprober/targets/rtc/rtcreporter"
)

const (
	sysvarsModuleName = "sysvars"
)

// Constants defining the default server host and port.
const (
	DefaultServerHost = ""
	DefaultServerPort = 9313
	ServerHostEnvVar  = "CLOUDPROBER_HOST"
	ServerPortEnvVar  = "CLOUDPROBER_PORT"
)

// Prober represents a collection of probes where each probe implements the Probe interface.
type Prober struct {
	Probes         map[string]*probes.ProbeInfo
	Servers        []*servers.ServerInfo
	c              *configpb.ProberConfig
	l              *logger.Logger
	rdsServer      *rdsserver.Server
	rtcReporter    *rtcreporter.Reporter
	Surfacers      []*surfacers.SurfacerInfo
	serverListener net.Listener

	// Used by GetConfig for /config handler.
	TextConfig string
}

func (pr *Prober) initDefaultServer() error {
	serverHost := pr.c.GetHost()
	if serverHost == "" {
		serverHost = DefaultServerHost
		// If ServerHostEnvVar is defined, it will override the default
		// server host.
		if host := os.Getenv(ServerHostEnvVar); host != "" {
			serverHost = host
		}
	}

	serverPort := int(pr.c.GetPort())
	if serverPort == 0 {
		serverPort = DefaultServerPort

		// If ServerPortEnvVar is defined, it will override the default
		// server port.
		if portStr := os.Getenv(ServerPortEnvVar); portStr != "" {
			port, err := strconv.ParseInt(portStr, 10, 32)
			if err != nil {
				return fmt.Errorf("failed to parse default port from the env var: %s=%s", ServerPortEnvVar, portStr)
			}
			serverPort = int(port)
		}
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", serverHost, serverPort))
	if err != nil {
		return fmt.Errorf("error while creating listener for default HTTP server. Err: %v", err)
	}
	pr.serverListener = ln

	return nil
}

// Init initialize prober with the given config file.
func (pr *Prober) Init(configFile string) error {
	// Initialize sysvars module
	l, err := logger.NewCloudproberLog(sysvarsModuleName)
	if err != nil {
		return err
	}
	sysvars.Init(l, nil)

	if pr.TextConfig, err = config.ParseTemplate(configFile, sysvars.Vars()); err != nil {
		return err
	}

	cfg := &configpb.ProberConfig{}
	if err := proto.UnmarshalText(pr.TextConfig, cfg); err != nil {
		return err
	}
	pr.c = cfg

	// Create a global logger. Each component gets its own logger on successful
	// creation. For everything else, we use a global logger.
	pr.l, err = logger.NewCloudproberLog("global")
	if err != nil {
		return fmt.Errorf("error in initializing global logger: %v", err)
	}

	// Initialize lameduck lister
	globalTargetsOpts := pr.c.GetGlobalTargetsOptions()

	if globalTargetsOpts.GetLameDuckOptions() != nil {
		ldLogger, err := logger.NewCloudproberLog("lame-duck")
		if err != nil {
			return fmt.Errorf("error in initializing lame-duck logger: %v", err)
		}

		if err := lameduck.InitDefaultLister(globalTargetsOpts.GetLameDuckOptions(), nil, ldLogger); err != nil {
			return err
		}
	}

	// Initiliaze probes
	pr.Probes, err = probes.Init(pr.c.GetProbe(), globalTargetsOpts, pr.l, sysvars.Vars())
	if err != nil {
		return err
	}

	// Start default HTTP server. It's used for profile handlers and
	// prometheus exporter.
	if err := pr.initDefaultServer(); err != nil {
		return err
	}

	// Initialize servers
	// TODO(manugarg): Plumb init context from cmd/cloudprober.
	initCtx, cancelFunc := context.WithCancel(context.TODO())
	pr.Servers, err = servers.Init(initCtx, pr.c.GetServer())
	if err != nil {
		cancelFunc()
		goto cleanupInit
	}

	pr.Surfacers, err = surfacers.Init(pr.c.GetSurfacer())
	if err != nil {
		goto cleanupInit
	}

	// Initialize RDS server, if configured.
	if c := pr.c.GetRdsServer(); c != nil {
		l, err := logger.NewCloudproberLog("rds-server")
		if err != nil {
			goto cleanupInit
		}
		if pr.rdsServer, err = rdsserver.New(initCtx, c, nil, l); err != nil {
			goto cleanupInit
		}
	}

	// Initialize RTC reporter, if configured.
	if opts := pr.c.GetRtcReportOptions(); opts != nil {
		l, err := logger.NewCloudproberLog("rtc-reporter")
		if err != nil {
			goto cleanupInit
		}
		if pr.rtcReporter, err = rtcreporter.New(opts, sysvars.Vars(), l); err != nil {
			goto cleanupInit
		}
	}
	return nil

cleanupInit:
	if pr.serverListener != nil {
		pr.serverListener.Close()
	}
	return err
}

// Start starts a previously initialized Cloudprober.
func (pr *Prober) Start(ctx context.Context) {
	// Start the default server
	srv := &http.Server{}
	go func() {
		<-ctx.Done()
		srv.Close()
	}()
	go func() {
		srv.Serve(pr.serverListener)
		os.Exit(1)
	}()

	dataChan := make(chan *metrics.EventMetrics, 1000)

	go func() {
		var em *metrics.EventMetrics
		for {
			em = <-dataChan
			var s = em.String()
			if len(s) > logger.MaxLogEntrySize {
				glog.Warningf("Metric entry for timestamp %v dropped due to large size: %d", em.Timestamp, len(s))
				continue
			}

			// Replicate the surfacer message to every surfacer we have
			// registered. Note that s.Write() is expected to be
			// non-blocking to avoid blocking of EventMetrics message
			// processing.
			for _, surfacer := range pr.Surfacers {
				surfacer.Write(context.Background(), em)
			}
		}
	}()

	// Start a goroutine to export system variables
	go sysvars.Start(ctx, dataChan, time.Millisecond*time.Duration(pr.c.GetSysvarsIntervalMsec()), pr.c.GetSysvarsEnvVar())

	// Start servers, each in its own goroutine
	for _, s := range pr.Servers {
		go s.Start(ctx, dataChan)
	}

	// Start RDS server if configured.
	if pr.rdsServer != nil {
		go pr.rdsServer.Start(ctx, dataChan)
	}

	// Start RTC reporter if configured.
	if pr.rtcReporter != nil {
		go pr.rtcReporter.Start(ctx)
	}

	if pr.c.GetDisableJitter() {
		for _, p := range pr.Probes {
			go p.Start(ctx, dataChan)
		}
		return
	}
	pr.startProbesWithJitter(ctx, dataChan)
}

// startProbesWithJitter try to space out probes over time, as much as possible,
// without making it too complicated. We arrange probes into interval buckets -
// all probes with the same interval will be part of the same bucket, and we
// then spread out probes within that interval by introducing a delay of
// interval / len(probes) between probes. We also introduce a random jitter
// between different interval buckets.
func (pr *Prober) startProbesWithJitter(ctx context.Context, dataChan chan *metrics.EventMetrics) {
	// Seed random number generator.
	rand.Seed(time.Now().UnixNano())

	// Make interval -> [probe1, probe2, probe3..] map
	intervalBuckets := make(map[time.Duration][]*probes.ProbeInfo)
	for _, p := range pr.Probes {
		intervalBuckets[p.Options.Interval] = append(intervalBuckets[p.Options.Interval], p)
	}

	for interval, probeInfos := range intervalBuckets {
		go func(interval time.Duration, probeInfos []*probes.ProbeInfo) {
			// Introduce a random jitter between interval buckets.
			randomDelayMsec := rand.Int63n(int64(interval.Seconds() * 1000))
			time.Sleep(time.Duration(randomDelayMsec) * time.Millisecond)

			interProbeDelay := interval / time.Duration(len(probeInfos))

			// Spread out probes evenly with an interval bucket.
			for _, p := range probeInfos {
				pr.l.Info("Starting probe: ", p.Name)
				go p.Start(ctx, dataChan)
				time.Sleep(interProbeDelay)
			}
		}(interval, probeInfos)
	}
}
