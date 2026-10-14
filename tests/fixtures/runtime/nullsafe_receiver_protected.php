<?php
class ProtectedBox { protected $secret = null; }
$o = new ProtectedBox;
echo $o?->secret?->take() ?? "fallback";
