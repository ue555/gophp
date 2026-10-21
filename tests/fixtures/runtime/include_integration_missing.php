<?php
$calls = 0;
function pathWithEffect() {
    global $calls;
    $calls++;
    return __DIR__ . '/../support/absent-integration.inc';
}
$first = include_once pathWithEffect();
$second = include_once pathWithEffect();
echo ($first ? 'yes' : 'no'), '|', ($second ? 'yes' : 'no'), '|', $calls, "\n";
