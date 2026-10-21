<?php
declare(strict_types=1);
echo __FILE__ === __DIR__ . '/include_integration_scope.inc' ? 'file|' : 'bad-file|';
$local++;
global $shared;
$shared++;
return $local;
