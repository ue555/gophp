<?php
function makeAppended(...$xs) {
    $xs[] = 30;
    return new Exception;
}
$args = makeAppended(10, 20)->getTrace()[0]['args'];
echo count($args), '|', $args[0], '|', $args[1], "\n";
