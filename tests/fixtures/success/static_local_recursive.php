<?php
function step($depth) {
    static $counter = 0;
    $counter++;
    if ($depth > 0) { step($depth - 1); }
    return $counter;
}
echo step(2), '|', step(0), "\n";
