<?php
function loadAutoThing($name) {
    if ($name === 'AutoThing') {
        require __DIR__ . '/../support/autoload_basic.inc';
    }
}
spl_autoload_register('loadAutoThing');
$thing = new AutoThing;
echo $thing->value(), '|';
echo spl_autoload_unregister('loadAutoThing') ? 'removed' : 'kept';
