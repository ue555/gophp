<?php
enum IntegratedMode: string { case On = 'on'; }
function integratedMode(BackedEnum $mode): string { return $mode->value; }
$callback = integratedMode(...);
echo $callback(IntegratedMode::On), '|';
echo call_user_func($callback, IntegratedMode::On), '|';
try { $callback('on'); } catch (TypeError $e) { echo 'type|'; }
echo IntegratedMode::from('on')->name;
