<?php
interface HasLabel { public function label(): string; }
trait LabelFromName { public function label(): string { return $this->name; } }
enum Signal implements HasLabel { use LabelFromName; case Stop; case Go; }
enum Status: string { case Open = 'o'; }
function identity(Signal $signal): UnitEnum { return $signal; }
function backed(Status $status): BackedEnum { return $status; }
class EnumSlot { public Signal|Status $item; }
$slot = new EnumSlot;
$slot->item = Signal::Stop;
echo identity(Signal::Go) === Signal::Go ? 'same|' : 'bad|';
echo Signal::Stop instanceof HasLabel ? 'interface|' : 'bad|';
echo Signal::Go->label(), '|';
echo backed(Status::Open)->value, '|', $slot->item->name;
