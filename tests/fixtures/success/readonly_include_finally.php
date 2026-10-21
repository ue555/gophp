<?php
require __DIR__ . '/../support/readonly_defs.inc';
$o = new IncludedReadonly(6);
try {
    $o->change();
} catch (Error $e) {
    echo 'caught|';
} finally {
    echo 'finally|';
}
echo $o->value, "\n";
