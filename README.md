# hamtraffic

## Positioning stations

To make the generated data feel less synthetic, stations from more populous locales appear more frequently. To
approximate the populousness, the following query was run in
[overpass turbo](https://overpass-turbo.eu/)
to get a
[GeoJSON](https://en.wikipedia.org/wiki/GeoJSON)
file listing all "cities" in
[OpenStreetMap](https://openstreetmap.org/copyright)
data:

```
[out:json];
(
node[place="city"];
);
out body;
>;
out skel qt;
```
