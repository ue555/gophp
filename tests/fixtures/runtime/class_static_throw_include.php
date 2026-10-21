<?php
echo 'before|';
try {
    require __DIR__ . '/../support/class_static_throw.inc';
} catch (Throwable $e) {
    echo 'caught|';
}
echo 'after|';
