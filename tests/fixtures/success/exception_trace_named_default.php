<?php
function makeNamedTrace($a = 1, $b = 2) { return new Exception('named'); }
$error = makeNamedTrace(b: 9);
$args = $error->getTrace()[0]['args'];
echo count($args), '|', $args[0], '|', $args[1], "\n";
