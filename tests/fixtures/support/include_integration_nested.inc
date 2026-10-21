<?php
try {
    require __DIR__ . '/include_integration_throw.inc';
} finally {
    echo 'inner-finally|';
}
