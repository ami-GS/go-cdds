/*
 * Copyright(c) 2006 to 2018 ADLINK Technology Limited and others
 *
 * This program and the accompanying materials are made available under the
 * terms of the Eclipse Public License v. 2.0 which is available at
 * http://www.eclipse.org/legal/epl-2.0, or the Eclipse Distribution License
 * v. 1.0 which is available at
 * http://www.eclipse.org/org/documents/edl-v10.php.
 *
 * SPDX-License-Identifier: EPL-2.0 OR BSD-3-Clause
 */
/****************************************************************

  Generated by Cyclone DDS IDL to C Translator
  File name: RoundTrip.h
  Source: RoundTrip.idl
  Generated: Mon Aug 27 23:30:42 PDT 2018
  Cyclone DDS: V0.1.0

*****************************************************************/

#include "ddsc/dds_public_impl.h"

#ifndef _DDSL_ROUNDTRIP_H_
#define _DDSL_ROUNDTRIP_H_


#ifdef __cplusplus
extern "C" {
#endif


typedef struct RoundTripModule_DataType
{
  dds_sequence_t payload;
} RoundTripModule_DataType;

extern const dds_topic_descriptor_t RoundTripModule_DataType_desc;

#define RoundTripModule_DataType__alloc() \
((RoundTripModule_DataType*) dds_alloc (sizeof (RoundTripModule_DataType)));

#define RoundTripModule_DataType_free(d,o) \
dds_sample_free ((d), &RoundTripModule_DataType_desc, (o))

#ifdef __cplusplus
}
#endif
#endif /* _DDSL_ROUNDTRIP_H_ */
