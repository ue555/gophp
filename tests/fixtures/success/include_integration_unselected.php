<?php
$calls = 0;
function pathWithEffect() {
    global $calls;
    $calls++;
    return __DIR__ . '/../support/absent-integration.inc';
}
if (false) {
    include pathWithEffect();
}
echo $calls, "\n";
