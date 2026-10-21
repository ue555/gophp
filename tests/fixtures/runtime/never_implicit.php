<?php
function broken(): never { echo "inside\n"; }
broken();
echo "after\n";
