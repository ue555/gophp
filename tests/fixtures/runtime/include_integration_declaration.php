<?php
echo 'before|';
try {
    require __DIR__ . '/../support/include_integration_duplicate.inc';
} catch (Throwable $e) {
    echo 'caught|';
}
echo 'after|';
