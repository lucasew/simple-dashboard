{ buildGoModule
, lib
}:

buildGoModule {
  pname = "simple-dashboard";
  version = "unstable-2022-07-15";

  src =  ./.;
  vendorHash = "sha256-a6iSGI+PJxIqF2WDp86SCR7Q2+pYf2kn0d7jKPScCyg=";

  postInstall = ''
      mkdir $out/share/simple-dashboard -p
      cp $src/*.ini* $out/share/simple-dashboard
  '';

  meta = with lib; {
    description = "Simple web-based dashboard to watch with your old tablet";
    homepage = "https://github.com/lucasew/simple-dashboard";
    license = licenses.mit;
    maintainers = with maintainers; [ lucasew ];
    mainProgram = "simple-dashboardd";
  };
}
