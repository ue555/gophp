<?php
function makeChanged(...$xs) {
    $xs[0] = 99;
    return new Exception;
}
$args = makeChanged(10, 20)->getTrace()[0]['args'];
echo count($args), '|', $args[0], '|', $args[1], "\n";
