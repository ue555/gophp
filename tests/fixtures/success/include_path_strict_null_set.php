<?php declare(strict_types=1);
$before = get_include_path();
try {
    set_include_path(null);
    echo 'accepted';
} catch (TypeError $e) {
    echo 'TypeError';
}
echo '|', get_include_path() === $before ? 'unchanged' : 'changed', "\n";
