package main

import "time"

func (this MonoSequencer) Next(slotCount int32) int64 {
	nextValue := this.pad.Load()
	nextSequence := nextValue + int64(slotCount)
	wrapPoint := nextSequence - int64(this.ringSize)
	cachedGatingSequence := this.pad[cachedGatingSequencePadIndex]

	if wrapPoint > cachedGatingSequence || cachedGatingSequence > nextValue {
		minSequence := int64(0)
		for wrapPoint > minSequence {
			minSequence = this.last.Load()
			time.Sleep(time.Nanosecond)
		}

		this.pad[cachedGatingSequencePadIndex] = minSequence
	}

	this.pad.Store(nextSequence)
	return nextSequence
}

func (this MonoSequencer) Publish(sequence int64) {
	this.cursor.Store(sequence)
}

func NewMonoSequencer(cursor *Sequence, ringSize int32, last Barrier) MonoSequencer {
	pad := NewSequence()
	pad[cachedGatingSequencePadIndex] = InitialSequenceValue

	return MonoSequencer{
		pad:      pad,
		cursor:   cursor,
		ringSize: ringSize,
		last:     last,
	}
}

type MonoSequencer struct {
	pad      *Sequence
	cursor   *Sequence
	ringSize int32
	last     Barrier
}

const cachedGatingSequencePadIndex = 1
