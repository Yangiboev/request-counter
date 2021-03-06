package server

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
)

type State struct {
	Past    RequestCounter
	Present Cache
}

type internalState struct {
	Past    requestCountList
	Present Cache
}

func (s State) encode() ([]byte, error) {
	internalState := internalState{
		Past:    s.Past.getNodes(),
		Present: s.Present,
	}
	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)
	err := e.Encode(internalState)
	if err != nil {
		return []byte{}, err
	}

	return b.Bytes(), nil
}

func decodeState(buffer []byte) (State, error) {
	var decodedInternalState internalState
	d := gob.NewDecoder(bytes.NewBuffer(buffer))
	err := d.Decode(&decodedInternalState)
	if err != nil {
		return State{}, err
	}

	decodedState := State{
		Past:    decodedInternalState.Past.ToRequestCounter(),
		Present: decodedInternalState.Present,
	}

	return decodedState, nil
}

func (s State) WriteToFile(path string) error {
	bytes, err := s.encode()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, bytes, 0600)
	if err != nil {
		return err
	}

	return nil
}

func ReadFromFile(path string) (State, error) {
	readBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return State{}, err
	}

	state, err := decodeState(readBytes)
	if err != nil {
		return State{}, err
	}

	return state, nil
}
