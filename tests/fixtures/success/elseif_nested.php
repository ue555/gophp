<?php
$x = 2;
IF ($x < 1) { echo "bad"; }
ELSEIF ($x == '2' && $x !== '2') {
    if ($x > 3) { echo "bad"; }
    else /* comment */ if ($x >= 2) {
        if (false) {} elseif (true) { echo "nested|"; }
    } else { echo "bad"; }
} else if (true) { echo "bad"; }
if (false) {} elseif (false) {} else if (true) {} else {}
echo "done\n";
