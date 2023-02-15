# hamtraffic

## World population data

The 30 arc-second world population
[GeoTIFF-format](https://en.wikipedia.org/wiki/GeoTIFF)
[data](https://hub.worldpop.org/geodata/summary?id=24777)
comes from the
[WorldPop Project](https://en.wikipedia.org/wiki/WorldPop_Project):

> WorldPop (www.worldpop.org - School of Geography and Environmental Science, University of Southampton; Department of Geography and Geosciences, University of Louisville; Departement de Geographie, Universite de Namur) and Center for International Earth Science Information Network (CIESIN), Columbia University (2018). Global High Resolution Population Denominators Project - Funded by The Bill and Melinda Gates Foundation (OPP1134076). https://dx.doi.org/10.5258/SOTON/WP00647

This large-ish
[data](./data/)
is decoupled from the repo with
[Git LFS](https://git-lfs.com/).

File details are as follows:

```console
$ listgeo data/ppp_2020_1km_Aggregated.tif
TIFFReadDirectory: Warning, Unknown field with tag 42113 (0xa481) encountered.
Geotiff_Information:
   Version: 1
   Key_Revision: 1.0
   Tagged_Information:
      ModelTiepointTag (2,3):
         0                 0                 0
         -180.001249265    83.99958319871    0
      ModelPixelScaleTag (1,3):
         0.0083333333      0.0083333333      0
      End_Of_Tags.
   Keyed_Information:
      GTModelTypeGeoKey (Short,1): ModelTypeGeographic
      GTRasterTypeGeoKey (Short,1): RasterPixelIsArea
      GeographicTypeGeoKey (Short,1): GCS_WGS_84
      GeogCitationGeoKey (Ascii,7): "WGS 84"
      GeogAngularUnitsGeoKey (Short,1): Angular_Degree
      GeogSemiMajorAxisGeoKey (Double,1): 6378137
      GeogInvFlatteningGeoKey (Double,1): 298.257223563
      End_Of_Keys.
   End_Of_Geotiff.

GCS: 4326/WGS 84
Datum: 6326/World Geodetic System 1984
Ellipsoid: 7030/WGS 84 (6378137.00,6356752.31)
Prime Meridian: 8901/Greenwich (0.000000/  0d 0' 0.00"E)
Projection Linear Units: User-Defined (1.000000m)

Corner Coordinates:
Upper Left    (180d 0' 4.50"W, 83d59'58.50"N)
Lower Left    (180d 0' 4.50"W, 72d 0' 1.50"S)
Upper Right   (179d59'55.50"E, 83d59'58.50"N)
Lower Right   (179d59'55.50"E, 72d 0' 1.50"S)
Center        (  0d 0' 4.50"W,  5d59'58.50"N)
```
