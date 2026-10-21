<?php
require __DIR__ . '/../support/void_never_defs.inc';
echo gettype(loadedVoid()), '|';
try { loadedNever(); } catch (Exception $e) { echo $e->getMessage(), "\n"; }
