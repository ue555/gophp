<?php
$counter = 1;
function bump() {
    global $counter;
    $counter += 2;
    return $counter;
}
echo bump(), '|', bump(), '|', $counter, "\n";
