<?php
$calls = 0;
function pathWithEffect() {
    global $calls;
    $calls++;
    return __DIR__ . '/../support/absent-integration.inc';
}
try {
    require_once pathWithEffect();
} catch (Error $e) {
    echo 'caught1|';
} finally {
    echo 'finally1|';
}
try {
    require_once pathWithEffect();
} catch (Error $e) {
    echo 'caught2|';
} finally {
    echo 'finally2|';
}
echo $calls, "\n";
