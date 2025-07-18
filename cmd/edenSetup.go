package cmd

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/lf-edge/eden/pkg/defaults"
	"github.com/lf-edge/eden/pkg/openevec"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newSetupCmd(configName, verbosity *string) *cobra.Command {
	cfg := &openevec.EdenSetupArgs{}
	var configDir, softSerial, zedControlURL, ipxeOverride string
	var grubOptions []string
	var netboot, installer bool

	var setupCmd = &cobra.Command{
		Use:               "setup",
		Short:             "setup harness",
		Long:              `Setup harness.`,
		PersistentPreRunE: preRunViperLoadFunction(cfg, configName, verbosity),
		Run: func(cmd *cobra.Command, args []string) {
			if err := openEVEC.SetupEden(*configName, configDir, softSerial, zedControlURL, ipxeOverride, grubOptions, netboot, installer); err != nil {

				log.Fatalf("Setup eden failed: %s", err)
			}
		},
	}

	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	setupCmd.Flags().BoolVarP(&cfg.Eden.Download, "download", "", cfg.Eden.Download, "download EVE or build")
	setupCmd.Flags().StringVar(&configDir, "eve-config-dir", filepath.Join(currentPath, "eve-config-dir"), "directory with files to put into EVE`s conf directory during setup")
	setupCmd.Flags().BoolVar(&netboot, "netboot", false, "Setup for use with network boot")
	setupCmd.Flags().BoolVar(&installer, "installer", false, "Setup for create installer")
	setupCmd.Flags().StringVar(&softSerial, "soft-serial", "", "Use provided serial instead of hardware one, please use chars and numbers here")
	setupCmd.Flags().StringVar(&zedControlURL, "zedcontrol", "", "Use provided zedcontrol domain instead of adam (as example: zedcloud.alpha.zededa.net)")
	setupCmd.Flags().StringVar(&ipxeOverride, "ipxe-override", "", "override lines inside ipxe, please use || as delimiter")
	setupCmd.Flags().StringArrayVar(&grubOptions, "grub-options", []string{}, "append lines to grub options")

	setupCmd.Flags().StringVarP(&cfg.Eden.CertsDir, "certs-dist", "o", cfg.Eden.CertsDir, "directory with certs")
	setupCmd.Flags().StringVarP(&cfg.Adam.CertsDomain, "domain", "d", defaults.DefaultDomain, "FQDN for certificates")
	setupCmd.Flags().StringVarP(&cfg.Adam.CertsIP, "ip", "i", defaults.DefaultIP, "IP address to use")
	setupCmd.Flags().StringVarP(&cfg.Adam.CertsEVEIP, "eve-ip", "", defaults.DefaultEVEIP, "IP address to use for EVE")
	setupCmd.Flags().StringVarP(&cfg.Eve.CertsUUID, "uuid", "u", defaults.DefaultUUID, "UUID to use for device")

	setupCmd.Flags().StringVarP(&cfg.Adam.Tag, "adam-tag", "", defaults.DefaultAdamTag, "Adam tag")
	setupCmd.Flags().StringVarP(&cfg.Adam.Dist, "adam-dist", "", cfg.Adam.Dist, "adam dist to start (required)")
	setupCmd.Flags().IntVarP(&cfg.Adam.Port, "adam-port", "", defaults.DefaultAdamPort, "adam dist to start")

	setupCmd.Flags().StringSliceVarP(&cfg.Eve.QemuFirmware, "eve-firmware", "", cfg.Eve.QemuFirmware, "firmware path")
	setupCmd.Flags().StringVarP(&cfg.Eve.QemuConfigPath, "config-path", "", cfg.Eve.QemuConfigPath, "path for config drive")
	setupCmd.Flags().StringVarP(&cfg.Eve.QemuDTBPath, "dtb-part", "", cfg.Eve.QemuDTBPath, "path for device tree drive (for arm)")
	setupCmd.Flags().StringVarP(&cfg.Eve.ImageFile, "image-file", "", cfg.Eve.ImageFile, "path for image drive (required)")
	setupCmd.Flags().StringVarP(&cfg.Eve.Dist, "eve-dist", "", cfg.Eve.Dist, "directory to save EVE")
	setupCmd.Flags().StringVarP(&cfg.Eve.Repo, "eve-repo", "", defaults.DefaultEveRepo, "EVE repo")
	setupCmd.Flags().StringVarP(&cfg.Eve.Registry, "eve-registry", "", defaults.DefaultEveRegistry, "EVE registry")
	setupCmd.Flags().StringVarP(&cfg.Eve.Tag, "eve-tag", "", defaults.DefaultEVETag, "EVE tag")
	setupCmd.Flags().StringVarP(&cfg.Eve.UefiTag, "eve-uefi-tag", "", defaults.DefaultEVETag, "EVE UEFI tag")
	setupCmd.Flags().StringVarP(&cfg.Eve.Arch, "eve-arch", "", runtime.GOARCH, "EVE arch")
	setupCmd.Flags().StringVarP(&cfg.Eve.Platform, "eve-platform", "", defaults.DefaultEVEPlatform, "EVE platform")
	setupCmd.Flags().StringToStringVarP(&cfg.Eve.HostFwd, "eve-hostfwd", "", defaults.DefaultQemuHostFwd, "port forward map")
	setupCmd.Flags().StringVarP(&cfg.Eve.QemuFileToSave, "qemu-config", "", cfg.Eve.QemuFileToSave, "file to save qemu config")
	setupCmd.Flags().StringVarP(&cfg.Eve.HV, "eve-hv", "", defaults.DefaultEVEHV, "hv of rootfs to use")

	setupCmd.Flags().StringVarP(&cfg.Eden.Images.EServerImageDist, "image-dist", "", cfg.Eden.Images.EServerImageDist, "image dist for eserver")
	setupCmd.Flags().StringVarP(&cfg.Eden.BinDir, "bin-dist", "", filepath.Join(currentPath, defaults.DefaultDist, defaults.DefaultBinDist), "directory for binaries")
	setupCmd.Flags().BoolVarP(&cfg.Adam.Force, "force", "", cfg.Adam.Force, "force overwrite config file")
	setupCmd.Flags().BoolVarP(&cfg.Adam.APIv1, "api-v1", "", cfg.Adam.APIv1, "use v1 api")

	setupCmd.Flags().StringVar(&cfg.Eve.BootstrapFile, "eve-bootstrap-file", "", "path to device config (in JSON) for bootstrapping")

	setupCmd.Flags().BoolVarP(&cfg.Eden.EnableIPv6, "enable-ipv6", "", false, "enable IPv6 connectivity for the Eden docker network")
	setupCmd.Flags().StringVarP(&cfg.Eden.IPv6Subnet, "ipv6-subnet", "", defaults.DefaultDockerNetIPv6Subnet, "IPv6 subnet for the Eden docker network")

	addSdnConfigDirOpt(setupCmd, cfg)
	addSdnImageOpt(setupCmd, cfg)
	addSdnDisableOpt(setupCmd, cfg)
	addSdnSourceDirOpt(setupCmd, cfg)
	addSdnLinuxkitOpt(setupCmd, cfg)

	return setupCmd
}
