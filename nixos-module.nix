{config, pkgs, lib, ...}:
let
  cfg = config.services.simple-dashboardd;
  module = pkgs.callPackage ./default.nix {};
in {
  options = with lib; {
    services.simple-dashboardd = {
      enable = mkEnableOption "Webapp to show system usage";
      config = mkOption {
        description = "Config string to be used by simple-dashboard";
        type = types.lines;
        default = builtins.readFile "${module}/share/simple-dashboard/config.ini.example";
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
      path = [ module ];
      script = ''
        ${module}/bin/simple-dashboardd -c ${builtins.toFile "simple-dashboard.cfg" cfg.config} -p ${builtins.toString cfg.port}
      '';
      restartIfChanged = true;
    };
    networking.firewall.allowedTCPPorts = mkIf cfg.openFirewall [ cfg.port ];
  };
}
