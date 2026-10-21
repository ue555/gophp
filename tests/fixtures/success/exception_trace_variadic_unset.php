<?php
function makeUnset(...$xs) {
    unset($xs[0]);
    return new Exception;
}
$args = makeUnset(10, 20)->getTrace()[0]['args'];
echo count($args), '|', $args[0], '|', $args[1], "\n";
