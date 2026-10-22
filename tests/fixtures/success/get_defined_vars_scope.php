<?php
$outer = 20;
function snapshot($argument) {
    global $outer;
    static $counter = 1;
    $local = [1, 2];
    $alias =& $argument;
    $removed = 3;
    unset($removed);
    $vars = get_defined_vars();
    ksort($vars);
    echo json_encode(array_keys($vars)), "\n";
    $vars['alias'] = 8;
    $vars['local'][] = 3;
    $vars['outer'] = 21;
    $vars['counter'] = 2;
    echo $argument, '|', count($local), '|', $vars['outer'], '|', $vars['counter'], "\n";
    echo $outer, '|', $counter, "\n";
}
snapshot(4);
class ViewScope {
    function render($data) {
        $vars = get_defined_vars();
        echo array_key_exists('this', $vars) ? "this included\n" : "this excluded\n";
        echo isset($vars['outer']) ? "leak\n" : "isolated\n";
    }
}
(new ViewScope)->render('value');
foreach ([[1], ['unknown' => 1]] as $args) {
    try { get_defined_vars(...$args); }
    catch (Throwable $e) { echo get_class($e), ': ', $e->getMessage(), "\n"; }
}
