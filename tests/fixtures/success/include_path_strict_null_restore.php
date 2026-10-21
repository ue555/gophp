<?php declare(strict_types=1);
set_include_path('temporary');
try {
    ini_restore(null);
    echo 'accepted';
} catch (TypeError $e) {
    echo 'TypeError';
}
echo '|', get_include_path() === 'temporary' ? 'unchanged' : 'changed', "\n";
