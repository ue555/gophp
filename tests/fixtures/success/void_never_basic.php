<?php
function done(): void { echo 'a'; return; }
function finish(): void { echo 'b'; }
function stop(): never { throw new Exception('stop'); }
echo gettype(done()), '|', gettype(finish()), '|';
try { stop(); } catch (Exception $e) { echo $e->getMessage(), "\n"; }
