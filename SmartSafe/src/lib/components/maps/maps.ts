import Map from "ol/Map";
import View from "ol/View";
import TileLayer from "ol/layer/Tile";
import XYZ from "ol/source/XYZ";
import { defaults as defaultControls, Zoom, Rotate } from "ol/control";
import Feature from "ol/Feature";
import Point from "ol/geom/Point";
import { Vector as VectorLayer } from "ol/layer";
import { Vector as VectorSource } from "ol/source";
import { Style } from "ol/style";
import { fromLonLat } from "ol/proj";
import CircleStyle from "ol/style/Circle";
import Fill from "ol/style/Fill";
import Stroke from "ol/style/Stroke";
import Control from "ol/control/Control";

function initMap(target: HTMLElement): Map {
  const userLocation = new Feature({ geometry: new Point([0, 0]) });
  let coords: [number, number] = [0, 0];

  class CenterToUserControl extends Control {
    constructor() {
      const div = document.createElement("div");
      div.innerHTML = `<button><i class="fa-solid fa-person-rays"></i>`;
      div.className = "center-to-user";
      div.addEventListener("click", () => {
        const view = map.getView();
        const zoomCurrent = view.getZoom();
        const zoomOut = view.getMinZoom() - 1;

        view.animate(
          {
            zoom: zoomOut,
            duration: 400,
            easing: (t) => t,
          },
          {
            center: coords,
            duration: 250,
            easing: (t) => t,
          },
          {
            zoom: zoomCurrent,
            duration: 400,
            easing: (t) => t,
          },
        );
      });

      super({
        element: div,
      });
    }
  }

  class EmergencyCall extends Control {
    constructor() {
      const div = document.createElement("div");
      div.innerHTML = `<button><i class="fa-solid fa-phone"></i></button>`;
      div.className = "emergency-call";
      div.addEventListener("click", () => {});

      super({
        element: div,
      });
    }
  }

  userLocation.setStyle(
    new Style({
      image: new CircleStyle({
        radius: 6,
        fill: new Fill({
          color: "#3399CC",
        }),
        stroke: new Stroke({
          color: "#fff",
          width: 2,
        }),
      }),
    }),
  );

  const userLayer = new VectorLayer({
    source: new VectorSource({
      features: [userLocation],
    }),
  });

  const map = new Map({
    target: target,
    layers: [
      new TileLayer({
        source: new XYZ({
          url: "https://{a-c}.tile.openstreetmap.org/{z}/{x}/{y}.png",
        }),
      }),
      userLayer,
    ],
    view: new View({
      center: [0, 0],
      zoom: 3,
      minZoom: 3,
    }),
    controls: defaultControls({ zoom: false, rotate: false }).extend([
      new Zoom({ className: "custom-zoom" }),
      new Rotate({ className: "custom-rotate" }),
      new CenterToUserControl(),
      new EmergencyCall(),
    ]),
  });

  let errorShown = false;
  let zoomToUser = false;

  function getUserLocation() {
    if (!navigator.geolocation) {
      console.error("Geolocation is not supported by your browser.");
      return;
    }

    navigator.geolocation.watchPosition(
      (position) => {
        const lon = position.coords.longitude;
        const lat = position.coords.latitude;
        coords = fromLonLat([lon, lat]) as [number, number];
        console.log("User Position:", coords);

        userLocation.setGeometry(new Point(coords));

        if (!zoomToUser) {
          zoomToUser = true;

          map.getView().setCenter(coords);
          map.getView().setZoom(15);
        }

        errorShown = false;
      },
      (error) => {
        if (!errorShown) {
          console.error("Error getting location:", error.message);
          errorShown = true;

          if (error.code === 1) {
            alert("Please enable location services in your browser settings.");
          } else if (error.code === 2) {
            alert("Location unavailable. Try again later.");
          } else if (error.code === 3) {
            alert("Location request timed out.");
          }
        }
      },
      {
        enableHighAccuracy: true,
        timeout: 10000,
        maximumAge: 0,
      },
    );
  }

  getUserLocation();

  return map;
}

export { initMap };
