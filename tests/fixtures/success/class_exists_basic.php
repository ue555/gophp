<?php
class KnownClass {}
interface KnownInterface {}
trait KnownTrait {}
echo class_exists('KnownClass') ? 'c|' : 'no|';
echo interface_exists('KnownInterface') ? 'i|' : 'no|';
echo trait_exists('KnownTrait') ? 't|' : 'no|';
echo class_exists('UnknownClass', false) ? 'wrong' : 'clear';
