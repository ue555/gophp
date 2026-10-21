<?php declare(strict_types=1);
require __DIR__ . '/../support/dnf_strict.inc';
try { crossTake('3'); } catch (TypeError $e) { echo 'TypeError'; }
echo '|', gettype(crossTake(3)), "\n";
