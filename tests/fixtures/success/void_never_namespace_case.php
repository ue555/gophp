<?php
namespace ReturnTypes;
function complete(): VoId { return; }
function stop(): NeVeR { throw new \Exception('stopped'); }
echo gettype(complete()), '|';
try { stop(); } catch (\Exception $e) { echo $e->getMessage(), "\n"; }
