# Galaxy
Porting some old JS code that created a procedurally generated galaxy.

Fairly huge number of stars (approximately Milky Way's worth) over roughly the Milky Way's volume, but random (without structure like arms etc.).

The portion of each class of stars is about correct.

## So Far
Got it showing a single sector that maps to the miniumum window.

Next add logc to handle viewport vs. sector misalignment.




Start by supporting more sectors along 1 axis in the viewport.
1. X Axis: Handle misalignment of sector with window. Cases are start (skip to start on viewport, insert all stars w/ offset, continue), stops (skip pre-beginning of viewport stars, - offset stars, finsh and stop short), continues (skip pre-beginning of viewport stars, finish stars, continue), completes (offset stars to viewport scaling down, finish stars, done).
2. X Axis: Add outer loop across sectors
3. X Axis: Handle window < sector, easier - as above but no completes case.
4. Repeat 1-3 for Y axis
5. Add rotation on X axis
6. Add rotation on T axis
7. Add rotation on Z axis
8. Add perspective
9. Add > 1 Z depth (as above, need to include fraction of Z above us in this sector plus additional sectors and a final ending sector, or smaller without middle or ven smallest wtihin < a single sector - but that can still crtoss sector boundaries) 
