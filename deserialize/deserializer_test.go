package deserializer

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/bnb-chain/zkbnb-setup/phase1"
	"github.com/bnb-chain/zkbnb-setup/phase2"
	"github.com/stretchr/testify/require"
)

///////////////////////////////////////////////////////////////////
///                             PTAU                            ///
///////////////////////////////////////////////////////////////////

// Format
// Taken from the iden3/snarkjs repo powersoftau_new.js file
// https://github.com/iden3/snarkjs/blob/master/src/powersoftau_new.js
/*
Header(1)
    n8
    prime
    power
tauG1(2)
    {(2 ** power)*2-1} [
        G1, tau*G1, tau^2 * G1, ....
    ]
tauG2(3)
    {2 ** power}[
        G2, tau*G2, tau^2 * G2, ...
    ]
alphaTauG1(4)
    {2 ** power}[
        alpha*G1, alpha*tau*G1, alpha*tau^2*G1,....
    ]
betaTauG1(5)
    {2 ** power} []
        beta*G1, beta*tau*G1, beta*tau^2*G1, ....
    ]
betaG2(6)
    {1}[
        beta*G2
    ]
contributions(7) - Ignore contributions, users can verify using snarkjs
    NContributions
    {NContributions}[
        tau*G1
        tau*G2
        alpha*G1
        beta*G1
        beta*G2
        pubKey
            tau_g1s
            tau_g1sx
            tau_g2spx
            alpha_g1s
            alpha_g1sx
            alpha_g1spx
            beta_g1s
            beta_g1sx
            beta_g1spx
        partialHash (216 bytes) See https://github.com/mafintosh/blake2b-wasm/blob/23bee06945806309977af802bc374727542617c7/blake2b.wat#L9
        hashNewChallenge
    ]
*/

func TestDeserializerPhase1(t *testing.T) {
	assert := require.New(t)

	input_path := "08.ptau"

	ptau, err := ReadPtau(input_path)

	if err != nil {
		assert.NoError(err)
	}

	phase1, err := convertPtauToSrs(ptau)

	fmt.Printf("TauG1: %v \n", phase1.Parameters.G1.Tau)
	fmt.Printf("AlphaTauG1: %v \n", phase1.Parameters.G1.AlphaTau)
	fmt.Printf("BetaTauG1: %v \n", phase1.Parameters.G1.BetaTau)
	fmt.Printf("TauG2: %v \n", phase1.Parameters.G2.Tau)
	fmt.Printf("BetaG2: %v \n", phase1.Parameters.G2.Beta)

	if err != nil {
		assert.NoError(err)
	}

	fmt.Printf("Size of the primes in bytes: %v \n", ptau.Header.n8)
}

func TestDeserializerPreparePhase2Ptau(t *testing.T) {
	assert := require.New(t)

	input_path := "08.ptau"

	ptau, err := ReadPtau(input_path)

	if err != nil {
		assert.NoError(err)
	}

	//mpcsetup.InitPhase2()

	fmt.Printf("Size of the primes in bytes: %v \n", ptau.Header.n8)
}

func TestDeserializePh1(t *testing.T) {
	assert := require.New(t)

	input_path_ptau := "08.ptau"

	ptau, err := ReadPtau(input_path_ptau)

	if err != nil {
		assert.NoError(err)
	}

	ph1, err := convertPtauToSrs(ptau)

	if err != nil {
		assert.NoError(err)
	}

	output_path := "08.ph1"

	outputFile, err := os.Create(output_path)

	if err != nil {
		assert.NoError(err)
	}

	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	var header phase1.Header

	header.Power = byte(uint8(8))

	// Taken from https://github.com/iden3/snarkjs/#7-prepare-phase-2
	header.Contributions = uint16(54)

	ph1.WriteTo(outputFile)

	phase1Path := "08.ph1"

	phase1File, err := os.Open(phase1Path)
	if err != nil {
		assert.NoError(err)
	}
	defer phase1File.Close()
}

func TestInitializePhase2(t *testing.T) {
	assert := require.New(t)

	ph1FilePath := "08.ph1"
	r1csFilePath := "demo_smtb.r1cs"
	phase2FilePath := "08.ph2"

	if err := phase2.Initialize(ph1FilePath, r1csFilePath, phase2FilePath); err != nil {
		assert.NoError(err)
	}

}

///////////////////////////////////////////////////////////////////
///                             ZKEY                            ///
///////////////////////////////////////////////////////////////////

// Taken from the iden3/snarkjs repo, zkey_utils.js
// (https://github.com/iden3/snarkjs/blob/fb144555d8ce4779ad79e707f269771c672a8fb7/src/zkey_utils.js#L20-L45)
// Format
// ======
// 4 bytes, zket
// 4 bytes, version
// 4 bytes, number of sections
// 4 bytes, section number
// 8 bytes, section size
// Header(1)
// 4 bytes, Prover Type 1 Groth
// HeaderGroth(2)
// 4 bytes, n8q
// n8q bytes, q
// 4 bytes, n8r
// n8r bytes, r
// 4 bytes, NVars
// 4 bytes, NPub
// 4 bytes, DomainSize  (multiple of 2)
//      alpha1
//      beta1
//      delta1
//      beta2
//      gamma2
//      delta2

func TestDeserializerZkey(t *testing.T) {
	input_path := "semaphore_16.zkey"

	assert := require.New(t)

	zkey, err := ReadZkey(input_path)

	if err != nil {
		assert.NoError(err)
	}

	fmt.Printf("ProtocolID for Groth16: %v \n", zkey.ZkeyHeader.ProtocolID)

	// protocolID should be 1 (Groth16)
	assert.Equal(GROTH_16_PROTOCOL_ID, zkey.ZkeyHeader.ProtocolID)

	fmt.Printf("n8q is: %v \n", zkey.protocolHeader.n8q)

	fmt.Printf("q is: %v \n", zkey.protocolHeader.q.String())

	fmt.Printf("n8r is: %v \n", zkey.protocolHeader.n8r)

	fmt.Printf("r is: %v \n", zkey.protocolHeader.r.String())

	fmt.Printf("nVars is: %v \n", zkey.protocolHeader.nVars)

	fmt.Printf("nPublic is: %v \n", zkey.protocolHeader.nPublic)

	fmt.Printf("domainSize is: %v \n", zkey.protocolHeader.domainSize)

	fmt.Printf("power is: %v \n", zkey.protocolHeader.power)
}
