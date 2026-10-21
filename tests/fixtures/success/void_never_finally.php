<?php
function complete(): void {
    try { throw new Exception('first'); }
    finally { echo 'finally|'; return; }
}
function replace(): never {
    try { throw new Exception('first'); }
    finally { throw new Exception('second'); }
}
complete();
try { replace(); } catch (Exception $e) { echo $e->getMessage(), "\n"; }
