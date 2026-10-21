<?php
function makeVariadicTrace(...$xs) { return new Exception('variadic'); }
$error = makeVariadicTrace(10, 20);
$args = $error->getTrace()[0]['args'];
echo count($args), '|', $args[0], '|', $args[1], "\n";
$named = makeVariadicTrace(label: 30);
$args = $named->getTrace()[0]['args'];
echo count($args), '|', $args['label'], "\n";
