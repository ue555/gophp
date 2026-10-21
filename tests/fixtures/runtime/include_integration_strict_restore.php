<?php
declare(strict_types=0);
function integrationAcceptInt(int $value): int { return $value; }
try {
    require __DIR__ . '/../support/include_integration_strict_throw.inc';
} catch (Exception $e) {
    echo $e->getMessage(), '|';
}
echo integrationAcceptInt(3.5), "\n";
