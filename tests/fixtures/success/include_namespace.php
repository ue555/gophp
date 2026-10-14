<?php
namespace Caller;
$loaded = require __DIR__ . '/../support/include_namespace.inc';
echo \Included\Space\answer(), '|', $loaded, "\n";
