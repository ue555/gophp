<?php
function outerCall() {
    $loads = 0;
    require __DIR__ . '/../support/include_integration_throw.inc';
}
echo 'before|';
outerCall();
echo 'after|';
