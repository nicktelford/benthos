// Copyright (c) 2018 Ashley Jeffs
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package processor

import (
	"github.com/Jeffail/benthos/lib/types"
)

//------------------------------------------------------------------------------

// ExecuteAll attempts to execute a slice of processors to a message. Returns
// N resulting messages or a response. The response may indicate either a NoAck
// in the event of the message being buffered or an unrecoverable error.
func ExecuteAll(procs []types.Processor, msgs ...types.Message) ([]types.Message, types.Response) {
	resultMsgs := make([]types.Message, len(msgs))
	copy(resultMsgs, msgs)

	var resultRes types.Response
	for i := 0; len(resultMsgs) > 0 && i < len(procs); i++ {
		var nextResultMsgs []types.Message
		for _, m := range resultMsgs {
			var rMsgs []types.Message
			if rMsgs, resultRes = procs[i].ProcessMessage(m); resultRes != nil && resultRes.Error() != nil {
				// We immediately return if a processor hits an unrecoverable
				// error on a message.
				return nil, resultRes
			}
			nextResultMsgs = append(nextResultMsgs, rMsgs...)
		}
		resultMsgs = nextResultMsgs
	}

	return resultMsgs, resultRes
}

//------------------------------------------------------------------------------

// FailFlagKey is a metadata key used for flagging processor errors in Benthos.
// If a message part has any non-empty value for this metadata key then it will
// be interpretted as having failed a processor step somewhere in the pipeline.
var FailFlagKey = "benthos_processing_failed"

// FlagFail marks a message part as having failed at a processing step.
func FlagFail(part types.Part) {
	part.Metadata().Set(FailFlagKey, "true")
}

// HasFailed checks whether a message part has failed a processing step.
func HasFailed(part types.Part) bool {
	return part.Metadata().Get(FailFlagKey) == "true"
}

// ClearFail removes any existing failure flags from a message part.
func ClearFail(part types.Part) {
	part.Metadata().Delete(FailFlagKey)
}

//------------------------------------------------------------------------------