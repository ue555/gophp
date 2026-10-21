<?php
try {
    echo require __DIR__ . '/../support/include_integration_finally_return.inc';
} catch (Exception $e) {
    echo 'wrong';
} finally {
    echo '|outer', "\n";
}
