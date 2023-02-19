# hamtraffic

## Configuration

All configuration takes place over environment variables. Here are the defaults:

```
FREERUN=true
RUNTIME= # n/a; this is parsed into a time.Duration
STATION_COUNT=10000
BANDS=160m:0.25,80m:0.40,40m:0.65,20m:1.0,10m:0.65,6m:0.40,2m:0.25
MODES=FT8:1.0,FT4:0.25,CW:0.15
STICKINESS=0.995
```

`BANDS` gives a relative weight for each, so that with the default setting above, it's four times more likely that a
station transmits on the 20-meter band, than on the 160-meter band, for example. The default also lists the bands that
the system currently understands. `MODES` follows the same pattern.

`STICKINESS` represents the likelihood of a station to stick to the band and mode it's currently working.

## Fake hams on air

To make it less likely to facepalm and get things mixed up with real hams, all callsigns are generated with prefix `X0`,
which [as of 2023-02-17 isn't allocated](https://en.wikipedia.org/wiki/Amateur_radio_call_signs).
The prefix is followed by a suffix of four letters from `AAAA` to `PPPP`. This range will provide around 65.5k unique
callsigns, which is enough for this purpose at this time.

### Positioning stations

To make the generated data feel less synthetic, stations get their locations from around 7.5k cities around the world.
The following query was run in
[overpass turbo](https://overpass-turbo.eu/)
to get a
[GeoJSON](https://en.wikipedia.org/wiki/GeoJSON)
file listing all cities in
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

It is a mystery to me why the query language looks like that (see
[here](https://wiki.openstreetmap.org/wiki/Overpass_API/FAQ#What_would_a_query_look_like_to_get_all_relations_tagged_with_type=boundary_or_type=multipolygon_and_their_way-members_and_the_nodes_used_by_those_way-members?)
for a more complex example), but I'm sure there are valid reasons.

The definition of what is considered a "city" seems to vary quite a bit in the data, currently the smallest one (Hmawbi)
has a population of 17. But anyway, that doesn't matter because the aim is to spread the stations around the world in an
"organic" pattern. Which doesn't mean that the station's locations would reflect in any way how amateur radio operators
are spread around the world, just that in those places there are (at least some) people. The main reason for taking this
approach is that I didn't want to deal with (massive) geography data to differentiate between land and water, for
example. Relatedly, maritime mobile isn't considered at all, yet.
