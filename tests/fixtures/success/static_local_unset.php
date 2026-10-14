<?php
function tick() {
    static $counter = 0;
    $counter++;
    echo $counter, '|';
    unset($counter);
    $counter = 99;
}
tick();
tick();
echo "done\n";
