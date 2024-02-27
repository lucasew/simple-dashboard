{config, pkgs, lib, ...}:
let
  cfg = config.services.simple-dashboardd;
in {
  options = with lib; {
    services.simple-dashboardd = {
      enable = mkEnableOption "Webapp to show system usage";
      package = mkOption {
        description = "Simple dashboard package";
        type = lib.types.package;
        default = pkgs.callPackage ./package.nix {};
      };
      config = mkOption {
        description = "Config string to be used by simple-dashboard";
        type = types.lines;
        default = builtins.readFile ./config.ini.example;
      };
      port = mkOption {
        description = "Port to listen";
        default = 42069;
        type = types.port;
      };
      openFirewall = mkEnableOption "Open the dashboard port on the firewall (recommended)";
    };
  };
  config = with lib; mkIf cfg.enable {
    systemd.services.simple-dashboardd = {
      enable = true;
      script = ''
        ${lib.getExe cfg.package} -c ${builtins.toFile "simple-dashboard.cfg" cfg.config} -p ${builtins.toString cfg.port}
      '';
      restartIfChanged = true;
    };
    networking.firewall.allowedTCPPorts = mkIf cfg.openFirewall [ cfg.port ];
  };
}
