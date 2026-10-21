<?php declare(strict_types=1);
$loaded = require __DIR__ . '/../support/integration_callable_contracts.inc';
$object = new IntegratedChild;
echo $loaded, '|', $object->label(), '|';
$callback = $object->callback();
try { $callback('2'); } catch (TypeError $e) { echo 'strict|'; }
echo $callback(2), '|';
echo call_user_func_array(args: ['n' => 4], callback: 'integratedNumber');
