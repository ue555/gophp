<?php
$error = require __DIR__ . '/../support/exception_trace_factory.inc';
$trace = $error->getTrace();
echo count($trace), '|', $trace[0]['function'], '|', $trace[1]['function'], '|';
echo $trace[0]['args'][0], '|', $error->getMessage(), "\n";
