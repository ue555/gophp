<?php
declare(strict_types=1);
try { new Exception(code: '3'); }
catch (TypeError $error) { echo $error->getMessage(), "\n"; }
