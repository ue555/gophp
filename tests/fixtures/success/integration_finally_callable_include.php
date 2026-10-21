<?php
function integratedRun(): void {
    try {
        echo require __DIR__ . '/../support/integration_thrower.inc', '|';
        $callback = integratedThrower();
        $callback();
    } finally {
        echo 'finally|';
    }
}
try { integratedRun(); } catch (Exception $e) { echo $e->getMessage(), '|'; }
echo 'after';
