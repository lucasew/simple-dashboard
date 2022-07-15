{pkgs ? import <nixpkgs> {}}:
pkgs.buildGoModule {
  pname = "simple-dashboard";
  version = "unstable-2022-07-15";
  src =  ./.;
  vendorSha256 = "sha256-a6iSGI+PJxIqF2WDp86SCR7Q2+pYf2kn0d7jKPScCyg=";
  postInstall = ''
      mkdir $out/share/simple-dashboard -p
      cp $src/*.ini* $out/share/simple-dashboard
  '';
  meta = with pkgs.lib; {
    description = "Simple web-based dashboard to watch with your old tablet";
    homepage = "https://github.com/lucasew/simple-dashboard";
    license = licenses.mit;
    maintainers = with maintainers; [ lucasew ];
  };
}
